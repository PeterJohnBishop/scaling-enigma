//
//  CreateUserView.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/20/25.
//

import Foundation
import SwiftUI

struct CreateUserView: View {
    @State var userVM: UserViewModel = UserViewModel()
    @State var password: String = ""
    @State var passwordConfirmation: String = ""
    @State var existingUser: Bool = false
    @State var showAlert: Bool = false
    @State var createSuccess: Bool = false
        
    var body: some View {
        NavigationStack{
                    VStack{
                        Spacer()
                        Text("Register").font(.system(size: 34))
                            .fontWeight(.ultraLight)
                        Divider().padding()
                        TextField("Name", text: $userVM.user.name)
                            .tint(.black)
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                            .padding()
                        TextField("Email", text: $userVM.user.email)
                            .tint(.black)
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                            .padding()
                        
                        SecureField("Password", text:  $password)
                            .tint(.black)
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                            .padding()
                        
                        
                        Button("Submit", action: {
                            userVM.user.password = password
                            Task{
                                let created = await userVM.createNewUser()
                                if created {
                                    createSuccess = await userVM.Login()
                                }
                            }
                        })
                        .fontWeight(.ultraLight)
                        .foregroundColor(.black)
                        .padding()
                        .background(
                            RoundedRectangle(cornerRadius: 8)
                                .fill(Color.white)
                                .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
                        )
                        .onChange(of: createSuccess, {
                            oldResponse, newResponse in
                            if !newResponse {
                                showAlert = true
                            }
                        })
                        .alert("Error", isPresented: $showAlert) {
                                        Button("OK", role: .cancel) {
                                            userVM.user.name = ""
                                            userVM.user.email = ""
                                            userVM.user.password = ""
                                        }
                                    } message: {
                                        Text(String(userVM.error ?? "Server Error!"))
                                    }
                                    .navigationDestination(isPresented: $createSuccess, destination: {
                                        SuccessView().navigationBarBackButtonHidden(true)
                                    })
                        Spacer()
                        HStack{
                            Spacer()
                            Text("I don't have an account.").fontWeight(.ultraLight)
                            Button("Login", action: {
                                existingUser = true
                            }).foregroundStyle(.black)
                                .fontWeight(.light)
                                .navigationDestination(isPresented: $existingUser, destination: {
                                    LoginUserView().navigationBarBackButtonHidden(true)
                                })
                            Spacer()
                        }
                    }.onAppear{
                        
                    }
                }
    }
}

#Preview {
    CreateUserView()
}
